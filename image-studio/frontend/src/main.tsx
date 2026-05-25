import React from 'react'
import {createRoot} from 'react-dom/client'
import './styles/index.css'
import App from './app/App'
import { applyPlatformAttributes } from './platform'
import { PlatformProvider } from './platform/context'
import './platform/android/wailsShim'

const container = document.getElementById('root')
applyPlatformAttributes()

const root = createRoot(container!)

root.render(
    <React.StrictMode>
        <PlatformProvider>
            <App/>
        </PlatformProvider>
    </React.StrictMode>
)
