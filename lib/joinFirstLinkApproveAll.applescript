set copyLinkShareButtonSize to {216, 22}

do shell script "open facetime://"
tell application "System Events"
	tell application process "FaceTime"
		tell first window -- default window
			tell scroll area 2 -- unlabelled
				tell (first list where description is "Recent Calls")
					tell (first list where description is "Upcoming")
						tell first group -- most recently created call/link
							repeat with anAction in (actions as list)
								if (description of anAction) is "FaceTime Video" then
									set desiredAction to anAction
									exit repeat
								end if
								--return anAction as action
							end repeat
							perform desiredAction
						end tell
					end tell
				end tell
			end tell
			tell button "Join"
				repeat until (exists)
					delay 0.1
				end repeat
				click
			end tell
			tell list 1 of list 1 of scroll area 2
				repeat until (exists)
					delay 0.1
				end repeat
				repeat with aUser in (every group as list)
					tell aUser
						repeat until (exists)
							delay 0.1
						end repeat
						repeat with checkImage in every image -- check if the element is a user row
							if description of checkImage is "contact silhouette" then
								tell (first button where description is "Approve join request")
									click
								end tell
							end if
						end repeat
					end tell
				end repeat
			end tell
		end tell
	end tell
end tell