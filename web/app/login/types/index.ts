export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginSuccessData {
  user_id: number;
  username: string;
  token: string;
}

export interface LoginSuccessResponse {
  status: "success";
  message: string;
  data: LoginSuccessData;
}

export interface LoginErrorResponse {
  status: "error";
  message: string;
}

export type LoginResponse = LoginSuccessResponse | LoginErrorResponse;

export interface LoginFormErrors {
  username?: string;
  password?: string;
}
