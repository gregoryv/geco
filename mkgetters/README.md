mkgetters

    $ mkgetters -t Car ./testdata/example.go
    // GENERATED!, DO NOT EDIT!
    
    package testdata
    
    func (c *Car) Model() string { return c.model }
    func (c *Car) Make() int     { return c.make }
