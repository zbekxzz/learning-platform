export interface UserPayload {
    email: string;
    full_name: string;
    role: 'student' | 'teacher' | 'admin';
    password?: string;
}
