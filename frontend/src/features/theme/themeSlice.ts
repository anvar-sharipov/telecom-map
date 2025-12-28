// import { createSlice } from '@reduxjs/toolkit';

// const initialState = {
//   theme: (localStorage.getItem('theme') as 'light' | 'dark') || 'light',
// };

// const themeSlice = createSlice({
//   name: 'theme',
//   initialState,
//   reducers: {
//     toggleTheme(state) {
//       state.theme = state.theme === 'light' ? 'dark' : 'light';
//     },
//   },
// });

// export const { toggleTheme } = themeSlice.actions;
// export default themeSlice.reducer;

import { createSlice } from '@reduxjs/toolkit';

type Theme = 'light' | 'dark';

const initialState = {
  theme: (localStorage.getItem('theme') as Theme) || 'light',
};

const themeSlice = createSlice({
  name: 'theme',
  initialState,
  reducers: {
    toggleTheme(state) {
      state.theme = state.theme === 'light' ? 'dark' : 'light';
      localStorage.setItem('theme', state.theme);
    },
  },
});

export const { toggleTheme } = themeSlice.actions;
export default themeSlice.reducer;
