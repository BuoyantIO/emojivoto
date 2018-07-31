import React from 'react';
import ReactDOM from 'react-dom';
import Vote from './components/Vote.jsx';
import Leaderboard from './components/Leaderboard.jsx';
import gridStyles from './../css/grid.css';
import styles from './../css/styles.css';

import { BrowserRouter, Route, Switch } from 'react-router-dom';

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
