# goxero

A simple go library for SquareUp Applications

## Call SquareUp Connect 
https://connect.squareup.com/oauth2/authorize?client_id={clientid}&scope={scope}

## Connect the client
v := gosquare.NewClient(code, clientId, clientSecret)

## Get an access toekn
token, merchant, expires, err := v.AccessToken()

or if you have a toekn

token, merchant, expires, err := v.RefreshToken(existingToken)

## Get Locations
locations, err := v.GetLocations(tooken)

## Setup Webhook
if err := v.UpdateWebHook(token, merchant, locationId, true); err != nil {
  service.HandleError(w, err, 422)
  return
}
