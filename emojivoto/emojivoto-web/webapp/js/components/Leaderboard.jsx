import React from 'react';
import _ from 'lodash';
import 'whatwg-fetch';
import { Link } from 'react-router-dom';


// TODO: Remove? (not currently)
// const randomLeaderboard = () => _.map(emo => {
//   return {
//     unicode: emo.unicode,
//     votes: Math.round(Math.random() * 1000)
//   }
// });
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

    // TODO: Remove? (not currently)
    // let leaderboard = randomLeaderboard();
    // this.setState({ leaderboard: _.orderBy(leaderboard, 'votes', 'desc') });
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
            <div className="btn btn-blue"><Link to="/">Vote on your favorite</Link></div>
            <div className="emoji-list">{this.renderLeaderboard()}
              <div className="footer-text">
                <p className="footer-experiment">A <a href='https://buoyant.io'>Buoyant</a> social experiment</p>
                <p>© 2017 Buoyant, Inc. All Rights Reserved.</p>
              </div>
            </div>
            </div>
            <div className="conduit-footer">
              <a href="https://conduit.io" target="_blank">
                <div className="footer-mark">
                <svg xmlns="http://www.w3.org/2000/svg">
                  <path fill="#ffffff" d="M26.23,25.27a1.13,1.13,0,0,0-.73.26,13.62,13.62,0,1,1,0-19.25h0a1.14,1.14,0,0,1-.07,1.51,1.15,1.15,0,0,1-1.63,0h0a11.34,11.34,0,1,0,.08,16.13,3.39,3.39,0,0,0,0-4.81,3.42,3.42,0,0,0-4.84,0,4.51,4.51,0,1,1,0-6.43h0a1.21,1.21,0,0,0,.72.26A1.14,1.14,0,0,0,21,11.78a1.17,1.17,0,0,0-.23-.72h0a6.83,6.83,0,1,0,0,9.66,1.12,1.12,0,0,1,1.6,0,1.15,1.15,0,0,1,0,1.62h0A9.07,9.07,0,1,1,22.21,9.4h0a3.41,3.41,0,0,0,4.91-4.72h0l-.06-.05,0,0h0a15.9,15.9,0,1,0,.08,22.52,1.07,1.07,0,0,0,.24-.71A1.12,1.12,0,0,0,26.23,25.27Z"/>
                </svg>
                </div>
                <div>
                  <p className="footer-cta">Tap here to learn more about Conduit</p>
                  <p className="footer-cta-web">Click here to learn more about Conduit</p>
                </div>
              </a>
            </div>
            </div>
            </div>
        </div>
    );
  }
}
