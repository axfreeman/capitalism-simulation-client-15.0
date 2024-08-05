## Various notes to self
* FetchRemoteTemplates creates a global template list. When and if users are allowed to create templates, this will have side effects

##Design
* hashed passwords should be recorded in the server: see `LoadRegisteredUsers`

##Bugs
* currently fails to display tables after taking an action (which it once did)  
* Does not resynchronise with the api server if a user login fails, which it should do.  


