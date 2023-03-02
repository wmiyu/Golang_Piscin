
1. Install ES from
    https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-7.17.0-darwin-x86_64.tar.gz
via:
    http://www.a.dukovany.cz/index.php?q=aHR0cHM6Ly93d3cuZWxhc3RpYy5jby9ndWlkZS9pbmRleC5odG1s

2. Oficial example
    https://github.com/elastic/go-elasticsearch.git

3. Test

curl -s -XGET "http://localhost:9200/places"

curl -s -XGET "http://localhost:9200/places/_doc/13000"

curl -s -XGET "http://localhost:9200/places/_search" 