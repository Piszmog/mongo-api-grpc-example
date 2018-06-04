package model

type Movie struct {
    Id          string `bson:"_id"`
    Name        string `bson:"name"`
    Description string `bson:"description"`
}

type CloudFoundryEnvironment struct {
    Mlab []MLab `json:"mlab"`
}
type MLab struct {
    Credentials Credentials `json:"credentials"`
}
type Credentials struct {
    Uri string `json:"uri"`
}
