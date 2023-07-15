mkgetters generates get methods for private struct fields

    $ mkgetters -h
    Usage: mkgetters [OPTIONS] [FILE]
    
    Options
        -t, --types : ""
            CSV list of types
    
        -w, --write-file : ""
        -h, --help

Example

    $ mkgetters -t Car ./testdata/example.go
    // GENERATED!, DO NOT EDIT!
    
    package testdata
    
    func (c *Car) Model() string { return c.model }
    func (c *Car) Make() int     { return c.make }
