# Mindia </!\\> Code under construction </!\\>

Mindia enables developers to manage their digital assets with their own storage. Mindia creates optimized assets with mutliple transformations defined by the user.

## Understanding Mindia

Mindia use a main storage to upload you assets into it. Transformed assets use second storage that acts like a cache. Mindia supports filesystem and S3 bucket storage.

To store api keys, named transformations, media metadatas, Mindia uses either the filesystem or an external storage. Mindia supports for now Redis, which can be used in our case for JSON document store, key/value store and search engine.

Mindia supports transformations on upload or when downloading media asset. An example here to create a thumbnail which we will be uploaded on the cache storage: "c_scale,h_200,w_300".

Next to transformations, analysis can also be performed on media assets. You can tag your pictures with Google Tagging. Tags will be stored on media's metadatas.

## Want to contribute ?

Want a new transformation, storage or third party integration ? Don't hesitate to contribute to
the project !
