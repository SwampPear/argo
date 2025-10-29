import './App.css'
import Content from './components/content/Content'
import SideMenu from './components/menus/side/Menu'
import TopMenu from './components/menus/top/Menu'


const App = () => {
    /*
    const updateResultText = (result: string) => setResultText(result);

    function greet() {
       Greet(name).then(updateResultText);
    }
    */

    return (
      <div id="app">
        <TopMenu/>
        <div id="container">
          <SideMenu/>
          <Content/>
        </div>
      </div>
    )
}

export default App
