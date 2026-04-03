export type UserRole = 'student' | 'teacher' | 'admin';

export interface UserProfile {
    email: string;
    full_name: string;
    id: number;
    role: UserRole;
}
