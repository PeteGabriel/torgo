# torgo

## Golang package to download from a magnet link

### Parse of magnet link

The file `magnet.go` defines the behavior used to parse and download from a given magnet link. First, a magnet link must begin with the suffix `magnet:?`. If not, the package does not acknowledge that the provided link is a valid magnet link. 

Every segment of the magnet link is separated by the character `&`. Each segment represent a key and a respective value in the form of `key=value`. At the moment the package searches for the keys represented in the list below.


1. "dn" - Display Name 
2. "xt" - Exact Topic 
3. "xs" - Exact Source 
4. "tr" - Tracker      
5. "ws" - Web Seed     
6. "xl" - Exact Length 

By getting each value for these keys we can build an instance of the type `Magnet` with the respective values. Later on this instace can be used to perform the actual download.


### Parse of .torrent file or uri

This package supports the parsing of a .torrent file or an uri of a .torrent file. It starts by searching if the file can be found locally. If not will try to download and save it locally. Once this is complete we can beging parsing the file. `.torrent` files use an encoding known as `bencode` and because of this I make use of an external package to help with the parsing of each segment. Once this is done an instance of `Torrent` is created and later on can be used to perform the actual download.
