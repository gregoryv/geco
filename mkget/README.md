mkget generates get methods for private struct fields

    $ mkget -h
    Usage: mkget [OPTIONS] [FILE]
    
    Options
        -t, --types : ""
            CSV list of types
    
        -w, --write-file : ""
        -a, --append-file : ""
        -h, --help

Example

    $ mkget -t Car
    // GENERATED, DO NOT EDIT!
    
    package main
    
    func (c *Car) Model() string { return c.model }
    func (c *Car) Make() int     { return c.make }
