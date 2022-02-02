package cache

import "github.com/keybase/go-keychain"

type keyChainOperation int

const (
	keyChainOperationGet = iota
	keyChainOperationPut
)

func (c *cacheImpl) GetSecret(service Service) ([]byte, bool) {
	query := c.createItem(service, keyChainOperationGet)
	results, err := keychain.QueryItem(query)
	if err != nil {
		panic(err)
	} else if len(results) != 1 {
		return nil, false
	}
	return results[0].Data, true
}

func (c *cacheImpl) PutSecret(service Service, secret []byte) {
	item := c.createItem(service, keyChainOperationPut)
	item.SetData(secret)
	if err := keychain.AddItem(item); err == keychain.ErrorDuplicateItem {
		prevItem := c.createItem(service, keyChainOperationGet)
		if results, err2 := keychain.QueryItem(prevItem); err2 != nil {
			panic(err2)
		} else if len(results) == 1 {
			if err3 := keychain.UpdateItem(prevItem, item); err3 != nil {
				panic(err3)
			}
		}
	} else if err != nil {
		panic(err)
	}
}

func (c *cacheImpl) createItem(service Service, op keyChainOperation) keychain.Item {
	item := keychain.NewItem()
	item.SetSecClass(keychain.SecClassGenericPassword)
	item.SetLabel(c.label)
	item.SetService(string(service))
	if op == keyChainOperationPut {
		item.SetSynchronizable(keychain.SynchronizableNo)
		item.SetAccessible(keychain.AccessibleWhenUnlocked)
	} else if op == keyChainOperationGet {
		item.SetReturnData(true)
	}

	return item
}
