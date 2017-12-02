import React from 'react';
import ReactDOM from 'react-dom';
import Vote from './components/Vote.jsx';
import Leaderboard from './components/Leaderboard.jsx';
import styles from './../css/styles.css';
import { BrowserRouter, Route, Switch } from 'react-router-dom';

// TODO: use static assets
// import logo from './../img/conduit-primary-white.svg';
// <img src={logo}/>

let appMain = document.getElementById('main');
let appData = appMain.dataset;

ReactDOM.render((
  <BrowserRouter>
    <div>
      <div className="main-content">
        <Switch>
          <Route exact path="/" component = { Vote } />
          <Route path="/leaderboard" component = { Leaderboard } />
        </Switch>
      </div>
    </div>
  </BrowserRouter>
), appMain);
