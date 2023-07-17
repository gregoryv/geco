mknozero generates method NoZero for given type which checks that each
private field is set to non zero value.

    $ mknozero -h
    Usage: mknozero [OPTIONS] [FILE]
    
    Options
        -t, --types : ""
            CSV list of types
    
        -w, --write-file : ""
        -a, --append-file : ""
        -h, --help

Example

    $ $ mknozero -t Car
    // GENERATED, DO NOT EDIT!
    
    package main
    
    // NoZero returns error if any private field is zero
    func (c *Car) NoZero() error {
            if reflect.ValueOf(c.model).IsZero() {
                    return fmt.Errorf("model not set")
            }
            if reflect.ValueOf(c.make).IsZero() {
                    return fmt.Errorf("make not set")
            }
            if reflect.ValueOf(c.output).IsZero() {
                    return fmt.Errorf("output not set")
            }
    }
