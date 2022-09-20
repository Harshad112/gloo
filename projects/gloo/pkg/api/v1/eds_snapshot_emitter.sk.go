// Code generated by solo-kit. DO NOT EDIT.

package v1

import (
	"sync"
	"time"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.uber.org/zap"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/errors"
	skstats "github.com/solo-io/solo-kit/pkg/stats"

	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/go-utils/errutils"
)

var (
	// Deprecated. See mEdsResourcesIn
	mEdsSnapshotIn = stats.Int64("eds.gloo.solo.io/emitter/snap_in", "Deprecated. Use eds.gloo.solo.io/emitter/resources_in. The number of snapshots in", "1")

	// metrics for emitter
	mEdsResourcesIn    = stats.Int64("eds.gloo.solo.io/emitter/resources_in", "The number of resource lists received on open watch channels", "1")
	mEdsSnapshotOut    = stats.Int64("eds.gloo.solo.io/emitter/snap_out", "The number of snapshots out", "1")
	mEdsSnapshotMissed = stats.Int64("eds.gloo.solo.io/emitter/snap_missed", "The number of snapshots missed", "1")

	// views for emitter
	// deprecated: see edsResourcesInView
	edssnapshotInView = &view.View{
		Name:        "eds.gloo.solo.io/emitter/snap_in",
		Measure:     mEdsSnapshotIn,
		Description: "Deprecated. Use eds.gloo.solo.io/emitter/resources_in. The number of snapshots updates coming in.",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{},
	}

	edsResourcesInView = &view.View{
		Name:        "eds.gloo.solo.io/emitter/resources_in",
		Measure:     mEdsResourcesIn,
		Description: "The number of resource lists received on open watch channels",
		Aggregation: view.Count(),
		TagKeys: []tag.Key{
			skstats.NamespaceKey,
			skstats.ResourceKey,
		},
	}
	edssnapshotOutView = &view.View{
		Name:        "eds.gloo.solo.io/emitter/snap_out",
		Measure:     mEdsSnapshotOut,
		Description: "The number of snapshots updates going out",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{},
	}
	edssnapshotMissedView = &view.View{
		Name:        "eds.gloo.solo.io/emitter/snap_missed",
		Measure:     mEdsSnapshotMissed,
		Description: "The number of snapshots updates going missed. this can happen in heavy load. missed snapshot will be re-tried after a second.",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{},
	}
)

func init() {
	view.Register(
		edssnapshotInView,
		edssnapshotOutView,
		edssnapshotMissedView,
		edsResourcesInView,
	)
}

type EdsSnapshotEmitter interface {
	Snapshots(watchNamespaces []string, opts clients.WatchOpts) (<-chan *EdsSnapshot, <-chan error, error)
}

type EdsEmitter interface {
	EdsSnapshotEmitter
	Register() error
	Upstream() UpstreamClient
}

func NewEdsEmitter(upstreamClient UpstreamClient, resourceNamespaceLister resources.ResourceNamespaceLister) EdsEmitter {
	return NewEdsEmitterWithEmit(upstreamClient, resourceNamespaceLister, make(chan struct{}))
}

func NewEdsEmitterWithEmit(upstreamClient UpstreamClient, resourceNamespaceLister resources.ResourceNamespaceLister, emit <-chan struct{}) EdsEmitter {
	return &edsEmitter{
		upstream:                upstreamClient,
		resourceNamespaceLister: resourceNamespaceLister,
		forceEmit:               emit,
	}
}

type edsEmitter struct {
	forceEmit <-chan struct{}
	upstream  UpstreamClient
	// resourceNamespaceLister is used to watch for new namespaces when they are created.
	// It is used when Expression Selector is in the Watch Opts set in Snapshot().
	resourceNamespaceLister resources.ResourceNamespaceLister
	// namespacesWatching is the set of namespaces that we are watching. This is helpful
	// when Expression Selector is set on the Watch Opts in Snapshot().
	namespacesWatching sync.Map
	// updateNamespaces is used to perform locks and unlocks when watches on namespaces are being updated/created
	updateNamespaces sync.Mutex
}

func (c *edsEmitter) Register() error {
	if err := c.upstream.Register(); err != nil {
		return err
	}
	return nil
}

func (c *edsEmitter) Upstream() UpstreamClient {
	return c.upstream
}

// Snapshots will return a channel that can be used to receive snapshots of the
// state of the resources it is watching
// when watching resources, you can set the watchNamespaces, and you can set the
// ExpressionSelector of the WatchOpts.  Setting watchNamespaces will watch for all resources
// that are in the specified namespaces. In addition if ExpressionSelector of the WatchOpts is
// set, then all namespaces that meet the label criteria of the ExpressionSelector will
// also be watched.
func (c *edsEmitter) Snapshots(watchNamespaces []string, opts clients.WatchOpts) (<-chan *EdsSnapshot, <-chan error, error) {

	if len(watchNamespaces) == 0 {
		watchNamespaces = []string{""}
	}

	for _, ns := range watchNamespaces {
		if ns == "" && len(watchNamespaces) > 1 {
			return nil, nil, errors.Errorf("the \"\" namespace is used to watch all namespaces. Snapshots can either be tracked for " +
				"specific namespaces or \"\" AllNamespaces, but not both.")
		}
	}

	errs := make(chan error)
	hasWatchedNamespaces := len(watchNamespaces) > 1 || (len(watchNamespaces) == 1 && watchNamespaces[0] != "")
	watchingLabeledNamespaces := !(opts.ExpressionSelector == "")
	var done sync.WaitGroup
	ctx := opts.Ctx

	// setting up the options for both listing and watching resources in namespaces
	watchedNamespacesListOptions := clients.ListOpts{Ctx: opts.Ctx, Selector: opts.Selector}
	watchedNamespacesWatchOptions := clients.WatchOpts{Ctx: opts.Ctx, Selector: opts.Selector}
	/* Create channel for Upstream */
	type upstreamListWithNamespace struct {
		list      UpstreamList
		namespace string
	}
	upstreamChan := make(chan upstreamListWithNamespace)
	var initialUpstreamList UpstreamList

	currentSnapshot := EdsSnapshot{}
	upstreamsByNamespace := sync.Map{}
	if hasWatchedNamespaces || !watchingLabeledNamespaces {
		// then watch all resources on watch Namespaces

		// watched namespaces
		for _, namespace := range watchNamespaces {
			/* Setup namespaced watch for Upstream */
			{
				upstreams, err := c.upstream.List(namespace, watchedNamespacesListOptions)
				if err != nil {
					return nil, nil, errors.Wrapf(err, "initial Upstream list")
				}
				initialUpstreamList = append(initialUpstreamList, upstreams...)
				upstreamsByNamespace.Store(namespace, upstreams)
			}
			upstreamNamespacesChan, upstreamErrs, err := c.upstream.Watch(namespace, watchedNamespacesWatchOptions)
			if err != nil {
				return nil, nil, errors.Wrapf(err, "starting Upstream watch")
			}

			done.Add(1)
			go func(namespace string) {
				defer done.Done()
				errutils.AggregateErrs(ctx, errs, upstreamErrs, namespace+"-upstreams")
			}(namespace)
			/* Watch for changes and update snapshot */
			go func(namespace string) {
				defer func() {
					c.namespacesWatching.Delete(namespace)
				}()
				c.namespacesWatching.Store(namespace, true)
				for {
					select {
					case <-ctx.Done():
						return
					case upstreamList, ok := <-upstreamNamespacesChan:
						if !ok {
							return
						}
						select {
						case <-ctx.Done():
							return
						case upstreamChan <- upstreamListWithNamespace{list: upstreamList, namespace: namespace}:
						}
					}
				}
			}(namespace)
		}
	}
	// watch all other namespaces that are labeled and fit the Expression Selector
	if opts.ExpressionSelector != "" {
		// watch resources of non-watched namespaces that fit the expression selectors
		namespaceListOptions := resources.ResourceNamespaceListOptions{
			Ctx:                opts.Ctx,
			ExpressionSelector: opts.ExpressionSelector,
		}
		namespaceWatchOptions := resources.ResourceNamespaceWatchOptions{
			Ctx:                opts.Ctx,
			ExpressionSelector: opts.ExpressionSelector,
		}

		filterNamespaces := resources.ResourceNamespaceList{}
		for _, ns := range watchNamespaces {
			// we do not want to filter out "" which equals all namespaces
			// the reason is because we will never create a watch on ""(all namespaces) because
			// doing so means we watch all resources regardless of namespace. Our intent is to
			// watch only certain namespaces.
			if ns != "" {
				filterNamespaces = append(filterNamespaces, resources.ResourceNamespace{Name: ns})
			}
		}
		namespacesResources, err := c.resourceNamespaceLister.GetResourceNamespaceList(namespaceListOptions, filterNamespaces)
		if err != nil {
			return nil, nil, err
		}
		newlyRegisteredNamespaces := make([]string, len(namespacesResources))
		// non watched namespaces that are labeled
		for i, resourceNamespace := range namespacesResources {
			c.namespacesWatching.Load(resourceNamespace)
			namespace := resourceNamespace.Name
			newlyRegisteredNamespaces[i] = namespace
			err = c.upstream.RegisterNamespace(namespace)
			if err != nil {
				return nil, nil, errors.Wrapf(err, "there was an error registering the namespace to the upstream")
			}
			/* Setup namespaced watch for Upstream */
			{
				upstreams, err := c.upstream.List(namespace, clients.ListOpts{Ctx: opts.Ctx})
				if err != nil {
					return nil, nil, errors.Wrapf(err, "initial Upstream list with new namespace")
				}
				initialUpstreamList = append(initialUpstreamList, upstreams...)
				upstreamsByNamespace.Store(namespace, upstreams)
			}
			upstreamNamespacesChan, upstreamErrs, err := c.upstream.Watch(namespace, clients.WatchOpts{Ctx: opts.Ctx})
			if err != nil {
				return nil, nil, errors.Wrapf(err, "starting Upstream watch")
			}

			done.Add(1)
			go func(namespace string) {
				defer done.Done()
				errutils.AggregateErrs(ctx, errs, upstreamErrs, namespace+"-upstreams")
			}(namespace)
			/* Watch for changes and update snapshot */
			go func(namespace string) {
				for {
					select {
					case <-ctx.Done():
						return
					case upstreamList, ok := <-upstreamNamespacesChan:
						if !ok {
							return
						}
						select {
						case <-ctx.Done():
							return
						case upstreamChan <- upstreamListWithNamespace{list: upstreamList, namespace: namespace}:
						}
					}
				}
			}(namespace)
		}
		if len(newlyRegisteredNamespaces) > 0 {
			contextutils.LoggerFrom(ctx).Infof("registered the new namespace %v", newlyRegisteredNamespaces)
		}

		// create watch on all namespaces, so that we can add all resources from new namespaces
		// we will be watching namespaces that meet the Expression Selector filter

		namespaceWatch, errsReceiver, err := c.resourceNamespaceLister.GetResourceNamespaceWatch(namespaceWatchOptions, filterNamespaces)
		if err != nil {
			return nil, nil, err
		}
		if errsReceiver != nil {
			go func() {
				for {
					select {
					case <-ctx.Done():
						return
					case err = <-errsReceiver:
						errs <- errors.Wrapf(err, "received error from watch on resource namespaces")
					}
				}
			}()
		}

		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case resourceNamespaces, ok := <-namespaceWatch:
					if !ok {
						return
					}
					// get the list of new namespaces, if there is a new namespace
					// get the list of resources from that namespace, and add
					// a watch for new resources created/deleted on that namespace
					c.updateNamespaces.Lock()

					// get the new namespaces, and get a map of the namespaces
					mapOfResourceNamespaces := make(map[string]bool, len(resourceNamespaces))
					newNamespaces := []string{}
					for _, ns := range resourceNamespaces {
						if _, hit := c.namespacesWatching.Load(ns.Name); !hit {
							newNamespaces = append(newNamespaces, ns.Name)
						}
						mapOfResourceNamespaces[ns.Name] = true
					}

					for _, ns := range watchNamespaces {
						mapOfResourceNamespaces[ns] = true
					}

					missingNamespaces := []string{}
					// use the map of namespace resources to find missing/deleted namespaces
					c.namespacesWatching.Range(func(key interface{}, value interface{}) bool {
						name := key.(string)
						if _, hit := mapOfResourceNamespaces[name]; !hit {
							missingNamespaces = append(missingNamespaces, name)
						}
						return true
					})

					for _, ns := range missingNamespaces {
						upstreamChan <- upstreamListWithNamespace{list: UpstreamList{}, namespace: ns}
					}

					for _, namespace := range newNamespaces {
						var err error
						err = c.upstream.RegisterNamespace(namespace)
						if err != nil {
							errs <- errors.Wrapf(err, "there was an error registering the namespace to the upstream")
							continue
						}
						/* Setup namespaced watch for Upstream for new namespace */
						{
							upstreams, err := c.upstream.List(namespace, clients.ListOpts{Ctx: opts.Ctx, Selector: opts.Selector})
							if err != nil {
								errs <- errors.Wrapf(err, "initial new namespace Upstream list in namespace watch")
								continue
							}
							upstreamsByNamespace.Store(namespace, upstreams)
						}
						upstreamNamespacesChan, upstreamErrs, err := c.upstream.Watch(namespace, clients.WatchOpts{Ctx: opts.Ctx, Selector: opts.Selector})
						if err != nil {
							errs <- errors.Wrapf(err, "starting new namespace Upstream watch")
							continue
						}

						done.Add(1)
						go func(namespace string) {
							defer done.Done()
							errutils.AggregateErrs(ctx, errs, upstreamErrs, namespace+"-new-namespace-upstreams")
						}(namespace)
						/* Watch for changes and update snapshot */
						go func(namespace string) {
							defer func() {
								c.namespacesWatching.Delete(namespace)
							}()
							c.namespacesWatching.Store(namespace, true)
							for {
								select {
								case <-ctx.Done():
									return
								case upstreamList, ok := <-upstreamNamespacesChan:
									if !ok {
										return
									}
									select {
									case <-ctx.Done():
										return
									case upstreamChan <- upstreamListWithNamespace{list: upstreamList, namespace: namespace}:
									}
								}
							}
						}(namespace)
					}
					if len(newNamespaces) > 0 {
						contextutils.LoggerFrom(ctx).Infof("registered the new namespace %v", newNamespaces)
					}
					c.updateNamespaces.Unlock()
				}
			}
		}()
	}
	/* Initialize snapshot for Upstreams */
	currentSnapshot.Upstreams = initialUpstreamList.Sort()

	snapshots := make(chan *EdsSnapshot)
	go func() {
		// sent initial snapshot to kick off the watch
		initialSnapshot := currentSnapshot.Clone()
		snapshots <- &initialSnapshot

		timer := time.NewTicker(time.Second * 1)
		previousHash, err := currentSnapshot.Hash(nil)
		if err != nil {
			contextutils.LoggerFrom(ctx).Panicw("error while hashing, this should never happen", zap.Error(err))
		}
		sync := func() {
			currentHash, err := currentSnapshot.Hash(nil)
			// this should never happen, so panic if it does
			if err != nil {
				contextutils.LoggerFrom(ctx).Panicw("error while hashing, this should never happen", zap.Error(err))
			}
			if previousHash == currentHash {
				return
			}

			sentSnapshot := currentSnapshot.Clone()
			select {
			case snapshots <- &sentSnapshot:
				stats.Record(ctx, mEdsSnapshotOut.M(1))
				previousHash = currentHash
			default:
				stats.Record(ctx, mEdsSnapshotMissed.M(1))
			}
		}

		defer func() {
			close(snapshots)
			// we must wait for done before closing the error chan,
			// to avoid sending on close channel.
			done.Wait()
			close(errs)
		}()
		for {
			record := func() { stats.Record(ctx, mEdsSnapshotIn.M(1)) }

			select {
			case <-timer.C:
				sync()
			case <-ctx.Done():
				return
			case <-c.forceEmit:
				sentSnapshot := currentSnapshot.Clone()
				snapshots <- &sentSnapshot
			case upstreamNamespacedList, ok := <-upstreamChan:
				if !ok {
					return
				}
				record()

				namespace := upstreamNamespacedList.namespace

				skstats.IncrementResourceCount(
					ctx,
					namespace,
					"upstream",
					mEdsResourcesIn,
				)

				// merge lists by namespace
				upstreamsByNamespace.Store(namespace, upstreamNamespacedList.list)
				var upstreamList UpstreamList
				upstreamsByNamespace.Range(func(key interface{}, value interface{}) bool {
					mocks := value.(UpstreamList)
					upstreamList = append(upstreamList, mocks...)
					return true
				})
				currentSnapshot.Upstreams = upstreamList.Sort()
			}
		}
	}()
	return snapshots, errs, nil
}
