mkenum generates enum method for comparable constants

    $ mkenum -h
    Usage: mkenum [OPTIONS] [FILE]
    
    Options
        -t, --types : ""
            CSV list of types
    
        -w, --write-file : ""
        -a, --append-file : ""
        -h, --help

Example

    $ mkenum -t Weekday
    // GENERATED, DO NOT EDIT!
    
    package main
    
    func (Weekday) Enum() []any { ... }
