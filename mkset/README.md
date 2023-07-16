mkset generates set methods for private struct fields

    $ mkset -h
    Usage: mkset [OPTIONS] [FILE]
    
    Options
        -t, --types : ""
            CSV list of types
    
        -w, --write-file : ""
        -a, --append-file : ""
        -h, --help

Example

    $ mkset -t Car .
    // GENERATED, DO NOT EDIT!
    
    package main
    
    func (c *Car) SetModel(v string) { c.model = v }
    func (c *Car) SetMake(v int)     { c.make = v }
