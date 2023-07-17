mkfill generates func Fill for copying values using set/get methods
for private struct fields

    $ mkfill
    Usage: mkfill [OPTIONS] [FILE]
    
    Options
        -t, --types : ""
            CSV list of types
    
        -w, --write-file : ""
        -a, --append-file : ""
        -h, --help

Example

    $ mkfill -t Car,Boat
    // GENERATED, DO NOT EDIT!
    
    package main
    
    func Fill(dst, src any) {
            {
                    dst, dstOk := dst.(interface{ SetModel(string) })
                    src, srcOk := src.(interface{ Model() string })
                    if dstOk && srcOk {
                            dst.SetModel(src.Model())
                    }
            }
            {
                    dst, dstOk := dst.(interface{ SetMake(int) })
                    src, srcOk := src.(interface{ Make() int })
                    if dstOk && srcOk {
                            dst.SetMake(src.Make())
                    }
            }
            {
                    dst, dstOk := dst.(interface{ SetOutput(io.Writer) })
                    src, srcOk := src.(interface{ Output() io.Writer })
                    if dstOk && srcOk {
                            dst.SetOutput(src.Output())
                    }
            }
            {
                    dst, dstOk := dst.(interface{ SetColor(int) })
                    src, srcOk := src.(interface{ Color() int })
                    if dstOk && srcOk {
                            dst.SetColor(src.Color())
                    }
            }
    }
