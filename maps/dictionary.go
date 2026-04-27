package main

import "errors"

type Dictionary map[string]string

func (d Dictionary) Search(word string) (string, error) {
	res, ok := d[word]

	if ok {
		return res, nil
	} else {
		return "", errors.New("Key not found")
	}
}

func (d Dictionary) Add(word, definition string) error {
	_, err := d.Search(word)

	if err == nil {
		return errors.New("Key already exists")
	}

	d[word] = definition
	return nil
}

func (d Dictionary) Update(word, definition string) error {
	_, err := d.Search(word)

	if err != nil {
		return errors.New("The key does not exist")
	}

	d[word] = definition
	return nil
}

func (d Dictionary) Delete(word string) error {
	_, err := d.Search(word)

	if err != nil {
		return errors.New("Key not found")
	}

	delete(d, word)
	return nil
}
