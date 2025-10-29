import './App.css'
import TopMenu from './components/menus/top/Menu'
import Page from './components/pages/Page'


const App = () => {
    return (
      <div id="app">
        <TopMenu/>
        <div id="container">
          <Page/>
        </div>
      </div>
    )
}

export default App
