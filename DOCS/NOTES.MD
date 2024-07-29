## Various notes to self
* FetchRemoteTemplates creates a global template list. When and if users are allowed to create templates, this will have side effects

##design
* apikey should not be repeated for each data item
* user.Sim is no longer needed
* hashed passwords should be recorded in the server: see `LoadRegisteredUsers`