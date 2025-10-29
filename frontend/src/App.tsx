import './App.css'
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
      <div id="App">
        <TopMenu/>
        <div id="Container">
          <SideMenu/>
        </div>
      </div>
    )
}

export default App
