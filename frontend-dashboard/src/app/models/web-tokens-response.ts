import { WebToken } from "./web-token";

export interface WebTokensResponse {
    status?: string;
    error_code?: string;
    data?: WebToken[];
}  