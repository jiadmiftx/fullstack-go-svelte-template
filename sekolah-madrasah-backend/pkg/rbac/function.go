package rbac

func getRouterDetails(find, parentPath string, routers []RouterRoles) []RouterDetail {
	var res []RouterDetail
	for _, router := range routers {
		currentPath := parentPath + router.Endpoint

		if find == currentPath {
			res = append(res, router.Detail...)
			return res
		}

		subRes := getRouterDetails(find, currentPath, router.SubRouter)
		res = append(res, subRes...)

		if len(subRes) > 0 {
			return res
		}
	}
	return res
}

func GetRouterDetails(find string, routers []RouterRoles) []RouterDetail {
	return getRouterDetails(find, "", routers)
}

func FindMany(routerRoles []RouterRoles, path string, filters ...Filter) (res []RouterDetail) {
	routers := GetRouterDetails(path, routerRoles)
	if len(routers) == 0 {
		return
	}
	if len(filters) == 0 {
		return routers
	}
	for _, filter := range filters {
		var currentLoopRes []RouterDetail
		if filter.Method != "" {
			var currentRes []RouterDetail
			for _, routerDetail := range routers {
				if routerDetail.Method == filter.Method {
					currentRes = append(currentRes, routerDetail)
				}
			}
			currentLoopRes = currentRes
		}

		if filter.Role != nil {
			var currentRes []RouterDetail
			currentData := routers
			if len(res) > 0 {
				currentData = res
			}
			for _, routerDetail := range currentData {
				for _, role := range routerDetail.Roles {
					if *filter.Role == role {
						currentRes = append(currentRes, routerDetail)
					}
				}
			}
			currentLoopRes = currentRes
		}
		res = append(res, currentLoopRes...)
	}
	return
}

func FindOne[T FindDataType](data []T, path string, filter *Filter) *RouterDetail {
	switch any(data).(type) {
	case []RouterRoles:
		currentFilter := Filter{}
		if filter != nil {
			currentFilter = *filter
		}
		dataConverted := any(data).([]RouterRoles)
		resMany := FindMany(dataConverted, path, currentFilter)
		if len(resMany) > 0 {
			return &resMany[0]
		}
		return nil
	case []RouterDetail:
		var res []RouterDetail
		routerDetailsConverted := any(data).([]RouterDetail)
		if len(routerDetailsConverted) == 0 {
			return nil
		}
		if filter != nil {
			if filter.Method != "" {
				var currentRes []RouterDetail
				for _, routerDetail := range routerDetailsConverted {
					if routerDetail.Method == filter.Method {
						currentRes = append(currentRes, routerDetail)
					}
				}
				res = currentRes
			}

			if filter.Role != nil {
				var currentRes []RouterDetail
				currentData := routerDetailsConverted
				if len(res) > 0 {
					currentData = res
				}
				for _, routerDetail := range currentData {
					for _, role := range routerDetail.Roles {
						if *filter.Role == role {
							currentRes = append(currentRes, routerDetail)
						}
					}
				}
				res = currentRes
			}
		}
		if len(res) != 0 {
			return &res[0]
		}
	}
	return nil
}
