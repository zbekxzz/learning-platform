export interface LoginResponse {
  token: string;
}

export interface RegisterRequest {
  full_name: string;
  email: string;
  password: string;
}
