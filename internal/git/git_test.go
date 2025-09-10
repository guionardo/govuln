package git

import "testing"

func TestParseGitURL(t *testing.T) {

	tests := []struct {
		name           string
		url            string
		wantHost       string
		wantOwner      string
		wantRepository string
		wantErr        bool
	}{
		{"invalid_name_should_return_error", "invalid_url", "", "", "", true},
		{"valid_git_should_return_data", "github.com/guionardo/govuln", "github.com", "guionardo", "govuln", false},
		{"config_ssh_git_should_return_data", "git@github.com:guionardo/govuln.git", "github.com", "guionardo", "govuln", false},
		{"config_https_git_should_return_data", "https://github.com/guionardo/govuln.git", "github.com", "guionardo", "govuln", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHost, gotOwner, gotRepository, err := ParseGitURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseGitURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotHost != tt.wantHost {
				t.Errorf("ParseGitURL() gotHost = %v, want %v", gotHost, tt.wantHost)
			}
			if gotOwner != tt.wantOwner {
				t.Errorf("ParseGitURL() gotOwner = %v, want %v", gotOwner, tt.wantOwner)
			}
			if gotRepository != tt.wantRepository {
				t.Errorf("ParseGitURL() gotRepository = %v, want %v", gotRepository, tt.wantRepository)
			}
		})
	}
}
