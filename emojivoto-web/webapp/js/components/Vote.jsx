import React from 'react';
import _ from 'lodash';
import { Link } from 'react-router-dom';
import 'whatwg-fetch';

const EmojiVotoPage = ({headline, contents, containerClass, preHeadline}) => {
  return (
    <div className={containerClass}>
      <div className="page-content container-fluid">
        <div className="row">
          <div className="col-md-12">
            {!preHeadline ? null : preHeadline}
            <h1 className="headline">{headline}</h1>

            {contents}
          </div>
        </div>
      </div>
    </div>
  );
}
export default class Vote extends React.Component {
  constructor(props) {
    super(props);
    this.resetState = this.resetState.bind(this);
    this.state = {
      emojiList: [],
      selectedEmoji: null,
      error: null
    }
  }

  loadFromServer() {
    fetch('/api/list').then(rsp => {
      rsp.json().then(r => {
        this.setState({ emojiList: r })
      }).catch(e => {
        this.setState( {error: e.statusText });
      }).catch(e => {
        this.setState( {error: e.statusText });
      });
    });
  }

  componentDidMount() {
    this.loadFromServer();
  }

  vote(emoji) {
    fetch(`/api/vote?choice=${emoji.shortcode}`).then(rsp => {
      if (rsp.ok) {
        this.setState({ selectedEmoji: emoji, error: null });
      } else {
        throw new Error("Unable to Register Vote");
      }
    }).catch(e => {
        this.setState({ selectedEmoji: emoji, error: e.toString() });
      });
  }

  resetState() {
    this.setState({ selectedEmoji: null, error: null });
  }

  renderEmojiList(emojis) {
    return _.map(emojis, (emoji, i) => {
      return (
        <div
          className="emoji emoji-votable"
          key={`emoji-${i}`}
          onClick={e => this.vote(emoji)}
        >
          {emoji.unicode}
        </div>
      );
    });
  }

  renderLeaderboardLink() {
    return <Link to="/leaderboard"><div className="btn btn-blue">View the leaderboard</div></Link>;
  }

  render() {
    if (this.state.error) {
      let errorMessage = <p>We couldn't process your request.</p>;
      if(this.state.selectedEmoji.shortcode === ":doughnut:") {
        errorMessage = (<div>
          <p className="doughnut-explanation">For the sake of this demo, voting for 🍩<br />
            always returns an error.
          </p>
          <p><Link to="/" onClick={this.resetState}>Pick another</Link>!</p>
        </div>);
      }

      let contents = (
        <div>
          {errorMessage}
          <Link to="/" onClick={this.resetState}><div className="btn btn-blue">Select again</div></Link>
        </div>
      );

      return <EmojiVotoPage
        preHeadline={<h1 className="title">404</h1>}
        headline="🚧"
        contents={contents}
        containerClass="background-500"
      />;
    } else if (!this.state.selectedEmoji) {
      let emojiList = this.state.emojiList;
      let contents = (
        <div>
          <h1>EMOJI VOTE</h1>
          <p>Tap to vote for your favorite emoji below</p>
          {this.renderLeaderboardLink()}
          {!_.isEmpty(emojiList) ? null : <div>Loading emoji...</div>}

          <div className="emoji-list">
            {this.renderEmojiList(emojiList)}

            <div className="footer-text">
              <p>A <a href='https://buoyant.io'>Buoyant</a> social experiment</p>
              <p>© 2017 Buoyant, Inc. All Rights Reserved.</p>
            </div>
          </div>
        </div>
      );

      return <EmojiVotoPage
        headline="🗳"
        contents={contents}
        containerClass="background"
      />;
    } else {
      let contents = (
        <div>
          <p>See how you stack up against others</p>
          {this.renderLeaderboardLink()}
          <Link to="/" onClick={this.resetState}><div className="btn btn-white">Pick another one</div></Link>
        </div>
      );
      return <EmojiVotoPage
        preHeadline={<h1>You picked:</h1>}
        headline={this.state.selectedEmoji.unicode}
        contents={contents}
        containerClass ="background"
      />;
    }
  }
}
