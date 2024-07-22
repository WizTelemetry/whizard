package store

/*
func (r *Store) horizontalPodAutoscalers() (retResources []resources.Resource) {
	timeRanges := r.store.Spec.TimeRanges
	if len(timeRanges) == 0 {
		timeRanges = append(timeRanges, v1alpha1.TimeRange{
			MinTime: r.store.Spec.MinTime,
			MaxTime: r.store.Spec.MaxTime,
		})
	}
	// for expected hpas
	var expectNames = make(map[string]struct{}, len(timeRanges))
	for i := range timeRanges {
		partitionSn := i
		partitionName := r.partitionName(i)
		expectNames[partitionName] = struct{}{}
		retResources = append(retResources, func() (runtime.Object, resources.Operation, error) {
			return r.horizontalPodAutoscaler(partitionName, partitionSn)
		})
	}

	var hpaList v2beta2.HorizontalPodAutoscalerList
	ls := r.BaseLabels()
	ls[constants.LabelNameAppName] = constants.AppNameStore
	ls[constants.LabelNameAppManagedBy] = r.store.Name
	err := r.Client.List(r.Context, &hpaList, client.InNamespace(r.store.Namespace), &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(ls),
	})
	if err != nil {
		return errResourcesFunc(err)
	}
	// check hpas to be deleted.
	for i := range hpaList.Items {
		hpa := hpaList.Items[i]
		if _, ok := expectNames[hpa.Name]; !ok {
			retResources = append(retResources, func() (runtime.Object, resources.Operation, error) {
				return &hpa, resources.OperationDelete, nil
			})
		}
	}
	return
}

func (r *Store) horizontalPodAutoscaler(name string, partitionSn int) (runtime.Object, resources.Operation, error) {
	var hpa = &v2beta2.HorizontalPodAutoscaler{ObjectMeta: r.meta(name, partitionSn)}

	if r.store.Spec.Scaler == nil {
		if err := r.Client.Get(r.Context, client.ObjectKeyFromObject(hpa), hpa); err != nil {
			if !util.IsNotFound(err) {
				return nil, "", err
			} else {
				return nil, "", nil
			}
		}
		// remove the existing hpa
		return hpa, resources.OperationDelete, nil
	}

	if err := r.Client.Get(r.Context, client.ObjectKeyFromObject(hpa), hpa); err != nil {
		if !util.IsNotFound(err) {
			return nil, "", err
		}
	}

	hpa.Spec.ScaleTargetRef = v2beta2.CrossVersionObjectReference{
		Kind:       "StatefulSet",
		APIVersion: "apps/v1",
		Name:       name,
	}

	if hpa.Labels == nil {
		hpa.Labels = r.labels(partitionSn)
	}

	hpa.Spec.MinReplicas = r.store.Spec.Scaler.MinReplicas
	hpa.Spec.MaxReplicas = r.store.Spec.Scaler.MaxReplicas
	hpa.Spec.Behavior = r.store.Spec.Scaler.Behavior
	hpa.Spec.Metrics = r.store.Spec.Scaler.Metrics

	return hpa, resources.OperationCreateOrUpdate, ctrl.SetControllerReference(r.store, hpa, r.Scheme)
}
*/
