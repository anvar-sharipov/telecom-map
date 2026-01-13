import { createSlice } from '@reduxjs/toolkit';
import type { PayloadAction } from '@reduxjs/toolkit';

export type NotificationType = 'success' | 'error' | 'info';

interface Notification {
  id: string;
  message: string;
  type: NotificationType;
}

interface NotificationState {
  list: Notification[];
}

const initialState: NotificationState = {
  list: [],
};

const notificationSlice = createSlice({
  name: 'notification',
  initialState,
  reducers: {
    showNotification: (state, action: PayloadAction<Omit<Notification, 'id'>>) => {
      state.list.push({
        id: crypto.randomUUID(),
        ...action.payload,
      });
    },
    removeNotification: (state, action: PayloadAction<string>) => {
      state.list = state.list.filter((n) => n.id !== action.payload);
    },
  },
});

export const { showNotification, removeNotification } = notificationSlice.actions;
export default notificationSlice.reducer;
