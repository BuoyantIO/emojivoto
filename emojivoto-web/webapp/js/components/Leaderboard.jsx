import React from 'react';
import _ from 'lodash';
import 'whatwg-fetch';
import { Link } from 'react-router-dom';

export default class Leaderboard extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      leaderboard: []
    }
  }

  componentDidMount() {
    this.loadFromServer();
  }

  loadFromServer(emoji) {
    fetch('/api/leaderboard').then(r => {
      r.json().then(emojis => {
        this.setState({
          leaderboard: emojis
        });
      }).catch(e => this.setState({ error: e }));
    }).catch(e => this.setState({ error: e }));
  }

  renderLeaderboard() {
    return _.map(this.state.leaderboard, (emoji, i) => {
      return (
        <div className="emoji" key={`emoji-${i}`} title={`${emoji.votes} votes`}>
          <div>{emoji.unicode}</div>
          { emoji.votes > 0 ? <div className="counter">{emoji.votes}</div> : null}
        </div>
      );
    });
  }

  render() {
    return (
      <div className="background">
        <div className="page-content container-fluid">
          <div className="row">
            <div className="col-md-12">
              {!this.state.error ? null : <div className="error">Error loading leaderboard.</div>}
              <h1>EMOJI VOTE LEADERBOARD </h1>
              <Link to="/"><div className="btn btn-blue">Vote on your favorite</div></Link>
              <div className="emoji-list">{this.renderLeaderboard()}
                <div className="footer-text">
                  <p className="footer-experiment">A <a href='https://buoyant.io'>Buoyant</a> social experiment</p>
                  <p>© 2018 Buoyant, Inc. All Rights Reserved.</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }
}
