import React from 'react'
import ReactDOM from 'react-dom/client'
import { RouterProvider } from 'react-router-dom'
import { Provider } from 'react-redux'
import { store } from './store/store'
import router from './routes'
import './index.css'
import 'antd/dist/reset.css'
// Removed AuthProvider import

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <Provider store={store}>
      {/* Removed AuthProvider wrapper */}
      <RouterProvider router={router} />
    </Provider>
  </React.StrictMode>,
)
