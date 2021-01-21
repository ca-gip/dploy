package ansible

import "github.com/ca-gip/dploy/internal/utils"

type Inventories []*Inventory

func (inventories Inventories) GetInventoryLength() int {
	return len(inventories)
}

// Returns variables name contained in the all section of all inventory files
func (inventories Inventories) GetInventoryKeys() (keys []string) {
	keySet := utils.NewSet()
	for _, inventory := range inventories {
		if inventory.Data != nil {
			for key := range inventory.Data.Groups["all"].Vars {
				keySet.Add(key)
			}
		}
	}
	return keySet.List()
}

// Given a variable name returns value across multiples inventory files
func (inventories Inventories) GetInventoryValues(key string) (values []string) {
	valueSet := utils.NewSet()
	for _, inventory := range inventories {
		if inventory.Data != nil {
			if value := inventory.Data.Groups["all"].Vars[key]; value != "" {
				valueSet.Add(value)
			}
		}
	}
	return valueSet.List()
}

// Returns inventory that match all the conditions
func (inventories Inventories) Filter(filters []Filter) (filtered []*Inventory) {
	if len(filters) == 0 {
		return
	}

	for _, inventory := range inventories {
		if inventory.Data != nil {

			type condition = string
			matchFilter := make(map[condition]bool)

			for _, filter := range filters {
				if filter.Eval(inventory.Data.Groups["all"].Vars[filter.Key]) {
					matchFilter[filter.GetRaw()] = true
				} else {
					matchFilter[filter.GetRaw()] = false
				}
			}

			if utils.MapHasAllTrue(matchFilter) {
				filtered = append(filtered, inventory)
			}

		}
	}
	return
}
