#webhooker - [![status](https://travis-ci.org/nkcraddock/webhooker.svg?branch=master)](https://travis-ci.org/nkcraddock/webhooker)
A simple webhooks implementation that uses rabbit for probably too much. (or not quite enough)

#For Subscribers
### Register as a webhooker

```
# Request
POST /api/register
{
    "url": "https://user-e-nelson.com/api/sourceapp-callback",
    "email": "user@email.com",
    "name": "Bropocalypse.com"
}

# Response
STATUS 201 Created
BODY "An api access key has been emailed to the specified address."

```

### Add a webhook
```
# Request
POST /api/webhook
{
    "evt": "user.*",
    "src": "sourceApp",
    "key": "Company_123"
}

# Response
STATUS 201 Created
HEADER "Location": "https://webhooker.com/api/webhook/c0bfa00b-02be-4493-9464-29f185836d4a"
```

### Wait for events!
Stand up an api listening at the url you specified when creating the webhook listening for messages such as:
```
# Request
POST https://user-e-nelson.com/api/sourceapp-callback
{
    "evt": "user.new",
    "src": "sourceApp",
    "key": "Company_123",
    "time": "2015-01-01 3:00:00 UCT", // or something
    "data": {
        "user_email": "somenewuser@email.com",
        "user_name": "Somenew User",
        "role": "king of the universe"
    }
}

# Response
STATUS 200 OK

```
#License
[MIT License](http://opensource.org/licenses/MIT)
