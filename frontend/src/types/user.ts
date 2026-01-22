import type { GroupDTO } from './group';

export type UserDTO = {
  id: number;
  username: string;
  full_name: string;
  is_active: boolean;
  created_at: string;
  groups: GroupDTO[];
};
