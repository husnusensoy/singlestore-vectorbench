Simple vector query benchmark tool to test single store capabilities.

## Run Vector Benchmark

Following command will run 100 SQL statements (`-c`) over prepared statement given by `--statement` by using randomly generated 512-dim vector using 
`-C` concurrent SingleStore sessions

```shell
bin/vectbench db -h 172.12.2.68,172.12.2.69 -p server -D deeplearning vbench \
-C 4 -c 100 \
--statement "select id, v<*> ? as distance from vecs order by distance use index(ivfpq_nlist) desc limit 100"
```



