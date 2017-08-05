# sequence

sequence is a package to print sequence files.

It is mainly for vfx industry, because they have so many of them.

It's simple package, so you could just add files and print it.


Basically it look like this.

```
man := sequence.NewManager(sequence.NewSplitter(), sequence.FmtSharp)
for _, f := range filenames {
	man.Add(f)
}

fmt.Println(man)
```

Please see example directory to see the full example.
