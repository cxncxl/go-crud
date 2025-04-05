package dto;

type CreateUserDto struct {
    Email    string `json:"email"`;
    Username string `json:"username"`;
    Password string `json:"password"`;
}

type PostLoginDto struct {
    Username string    `json:"username"`;
    Password string    `json:"password"`;
}
