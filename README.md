# Create automatic pull requests

## Prerequisites 


1. Setup your github account with a fork of the project
2. Fill a configuration json file like this :
   ```
   {
     "username": "github username",
     "name": "Firstname Lastname",
     "email": "email",
     "token": "github token here"
   }
   ```

## Usage

```
./pinata-auto-pr -organization myorg 
    -repository foobar \
    -file build.json \
    -key component.version \
    -value abcdef \
    -branch update-component \
    -message "Update component to abcdef" \
    -pr
```

What happens here ?
1. Retrieve build.json from myorg/foobar (master branch)
2. Replace ```component.version``` value by ```abcdef```
3. Create a branch ```update-component``` in username/foobar based on master
4. Update build.json
5. Create a PR in myorg/foobar

```-dry``` displays the new json file in stdout instead of creating a PR.
