# esm-archive

Archive multiple files into a single binary blob.

```js
import { Archive } from "esm-archive";

// bundle some files
const data = Archive.bundle([
  new File(["bar"], "foo.txt", { type: "text/plain" }),
  new File(["foo"], "bar.txt", { type: "text/plain" }),
]);

// use the archive
const archive = new Archive(data);
archive.checksum; // a 32-bit checksum of the archive
archive.entries[0].name; // => "foo.txt"
archive.entries[0].type; // => "text/plain"
archive.entries[1].name; // => "bar.txt"
archive.entries[1].type; // => "text/plain"
archive.readFile("foo.txt"); // => File(["bar"], "foo.txt", { type: "text/plain" })
archive.readFile("bar.txt"); // => File(["foo"], "bar.txt", { type: "text/plain" })
```

## Compression

This library does not compress the archived files. Here is an example of how to use `CompressionStream` to compress data
using gzip, and `DecompressionStream` to decompress it again.

```js
// bundle some files
const data = Archive.bundle([/* ... */]);

// compress the archive using CompressionStream
const compressed = await new Response(data).body.pipeThrough(new CompressionStream("gzip")).arrayBuffer();

// use the compressed archive
const decompressed = await new Response(compressed).body.pipeThrough(new DecompressionStream("gzip")).arrayBuffer();
const archive = new Archive(decompressed);
```

> Note that `CompressionStream` and `DecompressionStream` are not supported in all browsers, and you may need to use a
> polyfill or a different compression algorithm depending on your requirements.