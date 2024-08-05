# Frontend for the simulation of a capitalist economy  
* Interface with remote server via its API.  
* This version uses http://github.com/axfreeman/simulation-api-modified    
* There is more documentation there.
     
## Philosophy
The aim of this version is to focus on the UI logic by reducing the overhead of authentication to the bare minimum possible.  

The backend accepts only api-key header requests.  

This client uses an api-key issued by the server to access user data.  

This version uses only basic go http authentication.  

It creates a local database of registered users at startup and provides a facility to register a new user and synchronise that with the server.  

It does not yet resynchronise when a login fails, which it should do.
