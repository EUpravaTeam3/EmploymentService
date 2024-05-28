export class User{
    ID: string;
    firstName: string;
    lastName: string;
    ucn: string;
    password: string;
    email: string;
    address: string;
    type: string;
    token: string
    refreshToken: string
    constructor(){
        this.ID = ""
        this.firstName = ""
        this.lastName = ""
        this.ucn = ""
        this.password = ""
        this.email = ""
        this.address = ""
        this.type = ""
        this.token = ""
        this.refreshToken = ""
    }
}