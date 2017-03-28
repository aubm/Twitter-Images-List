## The images API

This module is a JSON/Rest API that is responsible for two things:

- Given an image absolute url, it will index it, among with a user provided description and a list a tags that are automatically computed on the fly.
  For exemple, given [this image](http://www.rd.com/wp-content/uploads/sites/2/2016/04/01-cat-wants-to-tell-you-laptop.jpg), the API will collect tags like "cat" and/or "laptop".
- Serving a search API with filter capabilities on the computed tags. So that one could use it to search for all cats with something like `/?tags=cat`.

The module is written in Go and is executed in Google App Engine standard. It uses [Google Datastore](https://cloud.google.com/datastore/docs/concepts/overview)
as a persistence backend, and the [Google Cloud Vision API](https://cloud.google.com/vision/) labels detection capability, to collect tags on the images.

## Usage

### Searching for images

```
GET https://images-api-dot-$PROJECT_ID.appspot.com/
```

Optional query parameters:

- `limit`: limits the max number of results to get in the response (default: 100).
- `offset`: shows the results starting from the requested offset (default: 0).
- `tags`: a list of tags to search for, seperated by blank spaces, example: "cat kitten".

Here is an example response:

```
{
   "items":[
      {
         "id":"923d7f5e-ebf2-4806-ab59-c031c561dea7",
         "url":"http://pbs.twimg.com/media/C31vlv6UMAECr1Y.jpg",
         "description":"RT @gossipgriII: facebook is so crazy https://t.co/RwT5Mrvl2B",
         "tags":[
            "text",
            "image",
            "font",
            "website",
            "screenshot",
            "brand",
            "presentation",
            "multimedia"
         ],
         "created_at":"2017-02-05T15:17:31.803496Z"
      },
      {
         "id":"acdfd801-f77f-4b94-bf21-cd4c020da51f",
         "url":"http://pbs.twimg.com/media/C36SWOAWYAUSkT0.jpg",
         "description":"RT @WhyLarryIsReal: listen https://t.co/hL3hAc79Oh",
         "tags":[
            "performance",
            "spring",
            "fashion",
            "sports"
         ],
         "created_at":"2017-02-05T15:17:31.905955Z"
      }
   ],
   "max_per_page":2,
   "offset":0,
   "count":100
}
```

### Indexing a new image

```
POST https://images-api-dot-$PROJECT_ID.appspot.com/index

{
  "url": "http://www.cats.org.uk/uploads/images/featurebox_sidebar_kids/grief-and-loss.jpg",
  "description": "Yet another cat"
}
```

### Queue the indexing of a new image

```
POST https://images-api-dot-$PROJECT_ID.appspot.com/queue-index

[
  {
    "url": "http://www.cats.org.uk/uploads/images/featurebox_sidebar_kids/grief-and-loss.jpg",
    "description": "Yet another cat"
  }
]
```

Here are the differences between `/index` and `/queue-index`:

- `/index` is synchronous and will process the provided image right away, while `/queue-index` is asynchonous.
  This means the service will provide you with an HTTP response and process the provided image later.
- `/index` can only be provided one image at a time, while `/queue-index` receives a list of indexing requests.

## Deploying the thing

### What you need

- The [Google Cloud SDK](https://cloud.google.com/sdk) in order to be able to use `gcloud` commands from your terminal. If you haven't, create a configuration with `gcloud init`.
- The [Google Cloud Vision API](https://console.cloud.google.com/apis/api/vision.googleapis.com/overview) enabled.
- An active Google API key. You can generate one [here](https://console.cloud.google.com/apis/credentials).
- To have the steps to create the `index-image` task queue [here](../gae-configs).

### Get it deployed

Go into the directory `images-api/app` and create a copy of the `app.yaml` file.

```
cp app.yaml app.prod.yaml
```

Then edit the file `app.prod.yaml`, add the environment bloc into it, then save and close it.

```
module: images-api
runtime: go
api_version: go1

env_variables:
  VISION_API_KEY: '{{CLOUD_VISION_API_KEY}}'

handlers:
- url: /.*
  script: _go_app
```

Replace the expression `{{CLOUD_VISION_API_KEY}}` by your Google Cloud Vision API key.

Deploy the module to App Engine.

```
gcloud app deploy --project $PROJECT_ID ./app.prod.yaml
```
