package main

const openFacetimeScript = `
do shell script "open facetime://"
delay 2
`

const checkFacetimeScript = `
tell application "System Events"
    set isRunning to (exists process "FaceTime")
end tell
return isRunning
`

const newLinkScript = `
set copyLinkShareButtonSize to {216, 22}

tell application "System Events"
	tell application process "FaceTime"
		tell first window
			repeat until (button "Create Link" exists)
				delay 0.1
			end repeat
			tell button "Create Link"
				click
				repeat until (first pop over exists)
					delay 0.1
				end repeat
				tell first pop over
					tell first group
						repeat until (third button exists)
							delay 0.1
						end repeat
						set childElements to every UI element
						repeat with childElement in childElements
							if childElement's size is copyLinkShareButtonSize then
								click childElement
								delay 0.1
								return (the clipboard as text)
								exit repeat
							end if
						end repeat
					end tell
				end tell
			end tell
		end tell
	end tell
end tell
`

const getActiveLinksScript = `
set linkList to {}

tell application "System Events"
	tell application process "FaceTime"
		tell first window -- default window
			tell second scroll area -- unlabelled
				tell (first list where description is "Recent Calls")
					if (first list where description is "Upcoming" exists)
						tell (first list where description is "Upcoming")
							repeat with linkEntry in every group
								tell linkEntry
									repeat with anAction in (actions as list)
										if (description of anAction) is "Info" then
											set desiredAction to anAction
											exit repeat
										end if
									end repeat
									perform desiredAction
									repeat until (first pop over exists)
										delay 0.1
									end repeat
									tell first pop over
										set linkList to (linkList & (name of (first static text where name contains "facetime.apple.com")))
									end tell
								end tell
							end repeat
						end tell
					end if
				end tell
			end tell
		end tell
	end tell
end tell
return linkList
`

const joinLatestLinkScript = `
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
		end tell
	end tell
end tell
`

const approveJoinScript = `
-- requires facetime to already be open, in a call, with the sidebar already toggled open
-- returns list of users that were admitted to the call
set admittedUsers to {}
tell application "System Events"
	tell first window of application process "FaceTime"
		tell scroll area 2
			tell list 1 -- description "Conversation Details"
				tell list 1 -- description "X people" where X is people already in call that aren't the host
					repeat with aPossibleUser in (every group as list)
						tell aPossibleUser
							set isUser to false
							tell (first button where description is "Approve join request")
								if (exists) then
									set isUser to true
									--click it
								end if
							end tell
							if isUser then
								tell static text 1
									set admittedUsers to admittedUsers & value of it
								end tell
							end if
						end tell
					end repeat
				end tell
			end tell
		end tell
	end tell
end tell
return admittedUsers
`

const joinLatestLinkAndApproveScript = `
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
							end repeat
							perform desiredAction
						end tell
					end tell
				end tell
			end tell
			delay 2 -- feels unnecessary be necessary but makes joining more reliable. Seems like the camera needs to initialize first
			tell button "Join"
				repeat until (exists)
					delay 0.1
				end repeat
				click
			end tell
			tell list 1 of list 1 of scroll area 2 -- poll since this needs to time to populate
				repeat until (exists)
					delay 0.1
				end repeat
				repeat with aPossibleUser in (every group as list)
					tell aPossibleUser
						repeat with checkImage in every image -- validate the element is a user row
							if description of checkImage is "contact silhouette" then
								tell (first button where description is "Approve join request") --this takes a beat to get displayed
									repeat until (exists)
										delay 0.1
									end repeat
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
`

const leaveCallScript = `
-- there is a button with description "End" for properly ending a call, but is hidden unless the sidebar is in focus.
tell application "System Events"
	tell first window of application process "FaceTime"
		tell attribute "AXCloseButton"
			click its value
		end tell
	end tell
end tell
`

const deleteLinkScript = `
set deleteTarget to "%s"
set foundMatch to false

tell application "System Events"
	tell application process "FaceTime"
		tell first window -- default window
			tell second scroll area -- unlabelled
				tell (first list where description is "Recent Calls")
					if ((first list where description is "Upcoming") exists) then
						tell (first list where description is "Upcoming")
							repeat with linkEntry in every group
								tell linkEntry
									repeat with anAction in (actions as list)
										if (description of anAction) is "Info" then
											set desiredAction to anAction
											exit repeat
										end if
									end repeat
									perform desiredAction
									repeat until (first pop over exists)
										delay 0.1
									end repeat
									tell first pop over
										repeat with checkText in (every static text as list)
											if (name of checkText is deleteTarget) then
												set foundMatch to true
												click button "Delete Link"
												exit repeat
											end if
										end repeat
									end tell
								end tell
							end repeat
						end tell
					end if
				end tell
			end tell
			if foundMatch then
				tell first sheet
					tell button "Delete Link"
						click
					end tell
				end tell
			end if
		end tell
	end tell
end tell
return foundMatch
`
