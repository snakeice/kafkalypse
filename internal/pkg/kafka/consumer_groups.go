package kafka

import "fmt"

func (k *Service) ConsumerGroups() ([]string, error) {
	cgs := make([]string, 0)

	groups, err := k.kAdmin.ListConsumerGroups()
	if err != nil {
		return nil, fmt.Errorf("failed to list consumer groups: %w", err)
	}

	for _, cg := range groups {
		cgs = append(cgs, cg)
	}

	return cgs, nil
}
