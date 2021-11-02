<h1 align="center">Gate God 🌉</h1>
<p>
  <a href="#" target="_blank">
    <img alt="License: MIT" src="https://img.shields.io/badge/License-MIT-yellow.svg" />
  </a>
</p>

> A project I built for my dad that automatically opens his house's front gate when I arrive.

Details on my blog: https://www.farley.ai/posts/gates

## Setup & Deployment

Firstly, create a free account with Plate Recognizer here: https://platerecognizer.com/ and copy your API token under /account. 

```sh
# Add API token to environment.
echo "PLATE_RECOGNIZER_TOKEN=<YOUR-API-TOKEN>" >> ./balena/.env-example && mv ./balena/.env-example ./balena/.env

# Build gate-god and move to our balena folder.
GOOS=linux GOARCH=arm GOARM=5 go build && mv ./gate-god ./balena/gate-god

# Deploy gate-god.
cd ./balena && balena push gate
```

Lastly, head over to balena/config/config.local.json and update the `allowed_plates` to all vehicle plates you wish to allow entry.
