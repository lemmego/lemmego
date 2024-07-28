package container

import (
	"fmt"
	"reflect"
	"sync"
)

// BindingType represents the type of binding
type BindingType int

const (
	// BindingTypeTransient represents a binding that's created each time it's resolved
	BindingTypeTransient BindingType = iota
	// BindingTypeSingleton represents a binding that's created only once
	BindingTypeSingleton
)

// binding represents a container binding
type binding struct {
	concrete interface{}
	bindType BindingType
	instance interface{}
}

// Container represents the IoC container
type Container struct {
	bindings map[reflect.Type]*binding
	mutex    sync.RWMutex
}

// New creates a new IoC container
func New() *Container {
	return &Container{
		bindings: make(map[reflect.Type]*binding),
	}
}

// Bind registers a transient binding with the container
func (c *Container) Bind(abstract interface{}, concrete interface{}) {
	c.bind(abstract, concrete, BindingTypeTransient)
}

// BindSingleton registers a singleton binding with the container
func (c *Container) BindSingleton(abstract interface{}, concrete interface{}) {
	c.bind(abstract, concrete, BindingTypeSingleton)
}

// bind is a helper method to register a binding
func (c *Container) bind(abstract interface{}, concrete interface{}, bindType BindingType) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	abstractType := c.getAbstractType(abstract)
	c.bindings[abstractType] = &binding{
		concrete: concrete,
		bindType: bindType,
	}
}

// Resolve retrieves a binding from the container
func (c *Container) Resolve(abstract interface{}) (interface{}, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	abstractType := c.getAbstractType(abstract)

	binding, exists := c.bindings[abstractType]
	if !exists {
		return nil, fmt.Errorf("binding not found for type: %v", abstractType)
	}

	if binding.bindType == BindingTypeSingleton && binding.instance != nil {
		return binding.instance, nil
	}

	instance, err := c.build(binding.concrete)
	if err != nil {
		return nil, err
	}

	if binding.bindType == BindingTypeSingleton {
		binding.instance = instance
	}

	return instance, nil
}

// Clear clears all bindings in the container
func (c *Container) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.bindings = make(map[reflect.Type]*binding)
}

// getAbstractType returns the reflect.Type of the abstract parameter
func (c *Container) getAbstractType(abstract interface{}) reflect.Type {
	if t, ok := abstract.(reflect.Type); ok {
		return t
	}
	return reflect.TypeOf(abstract)
}

// build creates an instance of the concrete type
func (c *Container) build(concrete interface{}) (interface{}, error) {
	t := reflect.TypeOf(concrete)

	if t.Kind() == reflect.Func {
		return c.buildFunc(concrete)
	}

	return concrete, nil
}

// buildFunc handles building instances from factory functions
func (c *Container) buildFunc(concrete interface{}) (interface{}, error) {
	t := reflect.TypeOf(concrete)
	params := make([]reflect.Value, t.NumIn())

	for i := 0; i < t.NumIn(); i++ {
		param := t.In(i)
		dependency, err := c.Resolve(param)
		if err != nil {
			return nil, fmt.Errorf("error resolving dependency: %v", err)
		}
		params[i] = reflect.ValueOf(dependency)
	}

	results := reflect.ValueOf(concrete).Call(params)
	if len(results) == 0 {
		return nil, nil
	}

	return results[0].Interface(), nil
}
