import { Tokens } from "./tokens";

export interface User {
    email?: string;
    last_login?: string;
    tokens?: Tokens;
    roles?: string
}